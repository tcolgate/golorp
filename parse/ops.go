// Copyright 2016 Tristan Colgate-McFarlane
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package parse

import "log"

type OpType int

const (
	XFX OpType = iota
	FX
	XFY
	FY
	YFX
	XF
	YF
)

type Op map[OpType]int

type OpSet map[string]Op

// defaultOps is a set of ops totally stolen from SWI, I have
// literally no idea what 90% of these do.
var defaultOps = OpSet{
	"-->":                   {XFX: 1200},
	":-":                    {XFX: 1200, FX: 1200},
	"?-":                    {FX: 1200},
	"dynamic":               {FX: 1150},
	"discontiguous":         {FX: 1150},
	"initialization":        {FX: 1150},
	"meta_predicate":        {FX: 1150},
	"module_transparent":    {FX: 1150},
	"multifile":             {FX: 1150},
	"public":                {FX: 1150},
	"thread_local":          {FX: 1150},
	"thread_initialization": {FX: 1150},
	"volatile":              {FX: 1150},
	";":                     {XFY: 1100},
	"|":                     {XFY: 1100},
	"->":                    {XFY: 1050},
	"*->":                   {XFY: 1050},
	",":                     {XFY: 1000},
	"*":                     {XFY: 1000, YFX: 400},
	":=":                    {XFX: 990},
	"\\+":                   {FY: 900},
	"<":                     {XFX: 700},
	"=":                     {XFX: 700},
	"=..":                   {XFX: 700},
	"=@=":                   {XFX: 700},
	"\\=@=":                 {XFX: 700},
	"=:=":                   {XFX: 700},
	"=<":                    {XFX: 700},
	"==":                    {XFX: 700},
	"=\\=":                  {XFX: 700},
	">":                     {XFX: 700},
	">=":                    {XFX: 700},
	"@<":                    {XFX: 700},
	"@=<":                   {XFX: 700},
	"@>":                    {XFX: 700},
	"@>=":                   {XFX: 700},
	"\\=":                   {XFX: 700},
	"\\==":                  {XFX: 700},
	"as":                    {XFX: 700},
	"is":                    {XFX: 700},
	">:<":                   {XFX: 700},
	":<":                    {XFX: 700},
	":":                     {XFY: 600},
	"+":                     {YFX: 500, FY: 200},
	"-":                     {YFX: 500, FY: 200},
	"/\\":                   {YFX: 500},
	"\\/":                   {YFX: 500},
	"xor":                   {YFX: 500},
	"?":                     {FX: 500},
	"/":                     {YFX: 400},
	"//":                    {YFX: 400},
	"div":                   {YFX: 400},
	"rdiv":                  {YFX: 400},
	"<<":                    {YFX: 400},
	">>":                    {YFX: 400},
	"mod":                   {YFX: 400},
	"rem":                   {YFX: 400},
	"**":                    {XFX: 200},
	"^":                     {XFY: 200},
	"\\":                    {FY: 200},
	".":                     {YFX: 100},
	"$":                     {FX: 1},
}

func (os OpSet) lookup(s string) (Op, bool) {
	ops, ok := os[s]
	if !ok {
		return nil, false
	}
	return ops, true
}

// Infix returns left, op, and right priorities
func (os OpSet) Infix(s string) (int, int, int, bool) {
	o, ok := os.lookup(s)
	log.Printf("check if %v is infix: %v %v\n", s, ok, o)
	if !ok {
		return 0, 0, 0, false
	}

	if opp, ok := o[YFX]; ok {
		return opp, opp, opp - 1, true
	}
	if opp, ok := o[XFY]; ok {
		return opp - 1, opp, opp, true
	}
	if opp, ok := o[XFX]; ok {
		return opp - 1, opp, opp - 1, true
	}
	return 0, 0, 0, false
}

func (os OpSet) Prefix(s string) (int, int, bool) {
	o, ok := os.lookup(s)
	log.Printf("check if %v is infix: \n", s, ok)
	log.Printf("check if %v is prefix: %v %v\n", s, ok, o)
	if !ok {
		return 0, 0, false
	}

	if opp, ok := o[FX]; ok {
		return opp, opp - 1, true
	}
	if opp, ok := o[FY]; ok {
		return opp, opp, true
	}
	return 0, 0, false
}

func (os OpSet) Postfix(s string) (int, int, bool) {
	o, ok := os.lookup(s)
	log.Printf("check if %v is postfix: %v %v\n", s, ok, o)
	if !ok {
		return 0, 0, false
	}

	if opp, ok := o[XF]; ok {
		return opp, opp - 1, true
	}
	if opp, ok := o[YF]; ok {
		return opp, opp, true
	}
	return 0, 0, false
}
