// Copyright 2018 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package wire

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"

	"golang.org/x/tools/go/types/typeutil"
)

type callKind int

const (
	funcProviderCall callKind = iota
	structProvider
	valueExpr
)

// A call represents a step of an injector function.  It may be either a
// function call or a composite struct literal, depending on the value
// of kind.
type call struct {
	// kind indicates the code pattern to use.
	kind callKind

	// out is the type this step produces.
	out types.Type

	// importPath and name identify the provider to call for kind ==
	// funcProviderCall or the type to construct for kind ==
	// structProvider.
	importPath string
	name       string

	// args is a list of arguments to call the provider with.  Each element is:
	// a) one of the givens (args[i] < len(given)), or
	// b) the result of a previous provider call (args[i] >= len(given))
	//
	// This will be nil for kind == valueExpr.
	args []int

	// fieldNames maps the arguments to struct field names.
	// This will only be set if kind == structProvider.
	fieldNames []string

	// ins is the list of types this call receives as arguments.
	// This will be nil for kind == valueExpr.
	ins []types.Type

	// The following are only set for kind == funcProviderCall:

	// hasCleanup is true if the provider call returns a cleanup function.
	hasCleanup bool
	// hasErr is true if the provider call returns an error.
	hasErr bool

	// The following are only set for kind == valueExpr:

	valueExpr     ast.Expr
	valueTypeInfo *types.Info
}

// solve finds the sequence of calls required to produce an output type
// with an optional set of provided inputs.
func solve(fset *token.FileSet, out types.Type, given []types.Type, set *ProviderSet) ([]call, error) {
	for i, g := range given {
		for _, h := range given[:i] {
			if types.Identical(g, h) {
				return nil, fmt.Errorf("multiple inputs of the same type %s", types.TypeString(g, nil))
			}
		}
	}
	providers, err := buildProviderMap(fset, set)
	if err != nil {
		return nil, err
	}

	// Start building the mapping of type to local variable of the given type.
	// The first len(given) local variables are the given types.
	index := new(typeutil.Map)
	for i, g := range given {
		if p := providers.At(g); p != nil {
			pp := p.(*Provider)
			return nil, fmt.Errorf("input of %s conflicts with provider %s at %s", types.TypeString(g, nil), pp.Name, fset.Position(pp.Pos))
		}
		index.Set(g, i)
	}

	// Topological sort of the directed graph defined by the providers
	// using a depth-first search. The graph may contain cycles, which
	// should trigger an error.
	var calls []call
	var visit func(trail []ProviderInput) error
	visit = func(trail []ProviderInput) error {
		typ := trail[len(trail)-1].Type
		if index.At(typ) != nil {
			return nil
		}
		for _, in := range trail[:len(trail)-1] {
			if types.Identical(typ, in.Type) {
				// TODO(light): Describe cycle.
				return fmt.Errorf("cycle for %s", types.TypeString(typ, nil))
			}
		}

		switch p := providers.At(typ).(type) {
		case nil:
			if len(trail) == 1 {
				return fmt.Errorf("no provider found for %s (output of injector)", types.TypeString(typ, nil))
			}
			// TODO(light): Give name of provider.
			return fmt.Errorf("no provider found for %s (required by provider of %s)", types.TypeString(typ, nil), types.TypeString(trail[len(trail)-2].Type, nil))
		case *Provider:
			if !types.Identical(p.Out, typ) {
				// Interface binding.  Don't create a call ourselves.
				if err := visit(append(trail, ProviderInput{Type: p.Out})); err != nil {
					return err
				}
				index.Set(typ, index.At(p.Out))
				return nil
			}
			for _, a := range p.Args {
				// TODO(light): This will discard grown trail arrays.
				if err := visit(append(trail, a)); err != nil {
					return err
				}
			}
			args := make([]int, len(p.Args))
			ins := make([]types.Type, len(p.Args))
			for i := range p.Args {
				ins[i] = p.Args[i].Type
				args[i] = index.At(p.Args[i].Type).(int)
			}
			index.Set(typ, len(given)+len(calls))
			kind := funcProviderCall
			if p.IsStruct {
				kind = structProvider
			}
			calls = append(calls, call{
				kind:       kind,
				importPath: p.ImportPath,
				name:       p.Name,
				args:       args,
				fieldNames: p.Fields,
				ins:        ins,
				out:        typ,
				hasCleanup: p.HasCleanup,
				hasErr:     p.HasErr,
			})
		case *Value:
			if !types.Identical(p.Out, typ) {
				// Interface binding.  Don't create a call ourselves.
				if err := visit(append(trail, ProviderInput{Type: p.Out})); err != nil {
					return err
				}
				index.Set(typ, index.At(p.Out))
				return nil
			}
			index.Set(typ, len(given)+len(calls))
			calls = append(calls, call{
				kind:          valueExpr,
				out:           typ,
				valueExpr:     p.expr,
				valueTypeInfo: p.info,
			})
		default:
			panic("unknown provider map value type")
		}
		return nil
	}
	if err := visit([]ProviderInput{{Type: out}}); err != nil {
		return nil, err
	}
	return calls, nil
}

func buildProviderMap(fset *token.FileSet, set *ProviderSet) (*typeutil.Map, error) {
	type binding struct {
		*IfaceBinding
		set *ProviderSet
	}

	providerMap := new(typeutil.Map) // to *Provider or *Value
	setMap := new(typeutil.Map)      // to *ProviderSet, for error messages
	var bindings []binding
	visited := make(map[*ProviderSet]struct{})
	next := []*ProviderSet{set}
	for len(next) > 0 {
		curr := next[0]
		copy(next, next[1:])
		next = next[:len(next)-1]
		if _, skip := visited[curr]; skip {
			continue
		}
		visited[curr] = struct{}{}
		for _, p := range curr.Providers {
			if providerMap.At(p.Out) != nil {
				return nil, bindingConflictError(fset, p.Pos, p.Out, setMap.At(p.Out).(*ProviderSet))
			}
			providerMap.Set(p.Out, p)
			setMap.Set(p.Out, curr)
		}
		for _, v := range curr.Values {
			if providerMap.At(v.Out) != nil {
				return nil, bindingConflictError(fset, v.Pos, v.Out, setMap.At(v.Out).(*ProviderSet))
			}
			providerMap.Set(v.Out, v)
			setMap.Set(v.Out, curr)
		}
		for _, b := range curr.Bindings {
			bindings = append(bindings, binding{
				IfaceBinding: b,
				set:          curr,
			})
		}
		for _, imp := range curr.Imports {
			next = append(next, imp)
		}
	}
	// Validate that bindings have their concrete type provided in the set.
	// TODO(light): Move this validation up into provider set creation.
	for _, b := range bindings {
		if providerMap.At(b.Iface) != nil {
			return nil, bindingConflictError(fset, b.Pos, b.Iface, setMap.At(b.Iface).(*ProviderSet))
		}
		concrete := providerMap.At(b.Provided)
		if concrete == nil {
			pos := fset.Position(b.Pos)
			typ := types.TypeString(b.Provided, nil)
			return nil, fmt.Errorf("%v: no binding for %s", pos, typ)
		}
		providerMap.Set(b.Iface, concrete)
		setMap.Set(b.Iface, b.set)
	}
	return providerMap, nil
}

// bindingConflictError creates a new error describing multiple bindings
// for the same output type.
func bindingConflictError(fset *token.FileSet, pos token.Pos, typ types.Type, prevSet *ProviderSet) error {
	position := fset.Position(pos)
	typString := types.TypeString(typ, nil)
	if prevSet.Name == "" {
		prevPosition := fset.Position(prevSet.Pos)
		return fmt.Errorf("%v: multiple bindings for %s (previous binding at %v)",
			position, typString, prevPosition)
	}
	return fmt.Errorf("%v: multiple bindings for %s (previous binding in %q.%s)",
		position, typString, prevSet.PkgPath, prevSet.Name)
}
