package validate

import (
	"github.com/graniticio/granitic/test"
	"github.com/graniticio/granitic/types"
	"testing"
)

func TestFloatTypeSupportDetection(t *testing.T) {

	fv := NewFloatValidatorBuilder("DEF", nil)

	sub := new(FloatsTarget)

	sub.F32 = 32.00
	sub.F64 = 64.1
	sub.NF = types.NewNilableFloat64(128E10)
	sub.S = "NAN"

	vc := new(validationContext)
	vc.Subject = sub

	checkFloatTypeSupport(t, "F64", vc, fv)
	checkFloatTypeSupport(t, "NF", vc, fv)
	checkFloatTypeSupport(t, "F32", vc, fv)

	bv, err := fv.parseRule("S", []string{"REQ:MISSING"})
	test.ExpectNil(t, err)

	_, err = bv.Validate(vc)

	test.ExpectNotNil(t, err)
}

func checkFloatTypeSupport(t *testing.T, it string, vc *validationContext, fvb *floatValidatorBuilder) {
	bv, err := fvb.parseRule(it, []string{"REQ:MISSING"})
	test.ExpectNil(t, err)

	r, err := bv.Validate(vc)
	test.ExpectNil(t, err)
	c := r.ErrorCodes

	test.ExpectInt(t, len(c), 0)
}

func TestFloatInSet(t *testing.T) {

	iv := NewFloatValidatorBuilder("DEF", nil)

	sub := new(FloatsTarget)

	sub.F64 = 3.0

	vc := new(validationContext)
	vc.Subject = sub

	bv, err := iv.parseRule("F64", []string{"REQ:MISSING", "IN:1,2,3,4,X"})
	test.ExpectNotNil(t, err)

	bv, err = iv.parseRule("F64", []string{"REQ:MISSING", "IN:1,2E10,3,4:NOT_IN"})
	test.ExpectNil(t, err)

	r, err := bv.Validate(vc)
	test.ExpectNil(t, err)
	c := r.ErrorCodes

	test.ExpectInt(t, len(c), 0)

	sub.F64 = 2.1E10

	r, err = bv.Validate(vc)
	test.ExpectNil(t, err)
	c = r.ErrorCodes

	test.ExpectInt(t, len(c), 1)
	test.ExpectString(t, c[0], "NOT_IN")

}

func TestFloatBreakOnError(t *testing.T) {

	iv := NewFloatValidatorBuilder("DEF", new(CompFinder))

	sub := new(FloatsTarget)

	sub.F64 = 3

	vc := new(validationContext)
	vc.Subject = sub

	bv, err := iv.parseRule("F64", []string{"REQ:MISSING", "BREAK"})
	test.ExpectNil(t, err)

	r, err := bv.Validate(vc)
	test.ExpectNil(t, err)
	c := r.ErrorCodes

	test.ExpectInt(t, len(c), 0)

	bv, err = iv.parseRule("F64", []string{"REQ:MISSING", "IN:1,2:NOTIN", "BREAK", "EXT:extFloat64Checker:EXTFAIL"})
	test.ExpectNil(t, err)

	r, err = bv.Validate(vc)
	test.ExpectNil(t, err)
	c = r.ErrorCodes

	test.ExpectInt(t, len(c), 1)

	test.ExpectString(t, c[0], "NOTIN")
}

func TestFloatRange(t *testing.T) {

	iv := NewFloatValidatorBuilder("DEF", nil)

	sub := new(FloatsTarget)

	sub.F32 = 3.1

	vc := new(validationContext)
	vc.Subject = sub

	bv, err := iv.parseRule("F32", []string{"REQ:MISSING", "RANGE:1-5"})
	test.ExpectNil(t, err)

	r, err := bv.Validate(vc)
	test.ExpectNil(t, err)
	c := r.ErrorCodes

	test.ExpectInt(t, len(c), 0)

	sub.F32 = 1.0

	r, err = bv.Validate(vc)
	test.ExpectNil(t, err)
	c = r.ErrorCodes

	test.ExpectInt(t, len(c), 0)

	sub.F32 = 5.0

	r, err = bv.Validate(vc)
	test.ExpectNil(t, err)
	c = r.ErrorCodes

	test.ExpectInt(t, len(c), 0)

	sub.F32 = -1.22

	r, err = bv.Validate(vc)
	test.ExpectNil(t, err)
	c = r.ErrorCodes

	test.ExpectInt(t, len(c), 1)

	sub.F32 = 6

	r, err = bv.Validate(vc)
	test.ExpectNil(t, err)
	c = r.ErrorCodes

	test.ExpectInt(t, len(c), 1)

	bv, err = iv.parseRule("F32", []string{"REQ:MISSING", "RANGE:-5"})
	sub.F32 = -20

	r, err = bv.Validate(vc)
	test.ExpectNil(t, err)
	c = r.ErrorCodes

	test.ExpectInt(t, len(c), 0)

	sub.F32 = 5

	r, err = bv.Validate(vc)
	test.ExpectNil(t, err)
	c = r.ErrorCodes

	test.ExpectInt(t, len(c), 0)

	sub.F32 = 6

	r, err = bv.Validate(vc)
	test.ExpectNil(t, err)
	c = r.ErrorCodes

	test.ExpectInt(t, len(c), 1)

	bv, err = iv.parseRule("F32", []string{"REQ:MISSING", "RANGE:5-"})
	sub.F32 = -20

	r, err = bv.Validate(vc)
	test.ExpectNil(t, err)
	c = r.ErrorCodes

	test.ExpectInt(t, len(c), 1)

	sub.F32 = 5

	r, err = bv.Validate(vc)
	test.ExpectNil(t, err)
	c = r.ErrorCodes

	test.ExpectInt(t, len(c), 0)

	sub.F32 = 6

	r, err = bv.Validate(vc)
	test.ExpectNil(t, err)
	c = r.ErrorCodes

	test.ExpectInt(t, len(c), 0)

}

type FloatsTarget struct {
	F32 float32
	F64 float64
	NF  *types.NillableFloat64
	S   string
}