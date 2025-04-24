package conv

import (
	"github.com/reeceappling/randomGoStuff/utils"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestOptions(t *testing.T) {
	t.Run("ensureNoDuplicates()", func(t *testing.T) {
		// TODO: this, only if used
	})

	t.Run("optionField", func(t *testing.T) {
		t.Run("OptionField", func(t *testing.T) {
			key, val := "aKey", 3728
			opt := OptionField(key, val)
			assert.Equal(t, key, opt.Key)
			assert.Equal(t, val, opt.val)
		})
	})

	t.Run("Options", func(t *testing.T) {
		var nilOpts *Options = nil
		testFields := []optionField{
			OptionField("a", "1"),
			OptionField("b", 2),
			OptionField("c", 3.0),
		}
		testOptsSingle := utils.Pointer(Options(map[string]any{testFields[0].Key: testFields[0].val}))
		t.Run("With", func(t *testing.T) {
			// Unnecessary with other tests
		})
		t.Run("Add", func(t *testing.T) {
			// Unnecessary with other tests
		})
		t.Run("Without", func(t *testing.T) {
			// Unnecessary with other tests
		})
		t.Run("Less", func(t *testing.T) {
			// Unnecessary with other tests
		})
		t.Run("addFields", func(t *testing.T) {
			t.Run("no fields", func(t *testing.T) {
				emptyFields := []optionField{}
				assert.Nil(t, nilOpts.addFields(false, emptyFields), "nil receiver should return nil")
				assert.True(t, reflect.DeepEqual(testOptsSingle, testOptsSingle.addFields(false, emptyFields)), "result should be as expected")
			})
			t.Run("with fields", func(t *testing.T) {
				opts := utils.Pointer(Options(map[string]any{testFields[0].Key: testFields[0].val}))
				res := opts.addFields(false, testFields)
				assert.NotNil(t, res)
				for _, field := range testFields {
					val, exists := (*res)[field.Key]
					assert.True(t, exists)
					assert.Equal(t, field.val, val)
				}
			})
		})
		t.Run("removeFields", func(t *testing.T) {
			assert.Nil(t, nilOpts.removeFields(false, []string{}), "nil receiver returns nil")
			t.Run("empty output returns nil", func(t *testing.T) {
				assert.Nil(t, testOptsSingle.removeFields(true, []string{testFields[0].Key}))
			})
			t.Run("non-empty output returns non-nil", func(t *testing.T) {
				assert.True(t, reflect.DeepEqual(testOptsSingle, testOptsSingle.removeFields(false, []string{})))
			})
		})
		t.Run("cloneIfNeeded", func(t *testing.T) {
			original := utils.Pointer(Options(map[string]any{}))
			opts := original.cloneIfNeeded(false)
			assert.Equal(t, original, opts)
			newOpts := *opts
			newOpts[testFields[0].Key] = testFields[0].val
			*opts = newOpts
			assert.False(t, reflect.DeepEqual(*original, *opts), "modifying result should modify original receiver for false case")
		})
		t.Run("cloneIfNeededPossiblyNil", func(t *testing.T) {
			// TODO: test that nil options returns a pointer to an empty options
			// TODO: test non-nil options returns a pointer to an equivalent underlying
		})
		t.Run("Clone", func(t *testing.T) {
			assert.Nil(t, nilOpts.Clone(), "nil case")
			// TODO: ensure modifying result doesn't modify original
		})
	})

	t.Run("CombineOptions()", func(t *testing.T) {
		// TODO: no options returns nil
		// TODO: all nil returns nil
		// TODO: all but one nil returns the one
		// TODO: value collision (same value, different value)
		// TODO: combines as expected
	})

	//t.Run("Options", func(t *testing.T) {
	//	blankOpts := NewOptions()
	//	t.Run("NewOptions()", func(t *testing.T) {
	//		t.Run("no options returns nil", func(t *testing.T) {
	//			assert.Nil(t, blankOpts)
	//		})
	//		t.Run("with options returns correctly", func(t *testing.T) {
	//			// TODO: this
	//		})
	//	})
	//	keyA, keyB := "keyA", "keyB"
	//	valA, valB := 3728, "notAnInt"
	//	fields := []optionField{
	//		OptionField(keyA, valA),
	//		OptionField(keyB, valB),
	//	}
	//	opts := blankOpts.With(fields...)
	//	t.Run("With()", func(t *testing.T) {
	//		t.Run("Adding fields", func(t *testing.T) {
	//			assert.Equal(t, len(fields), len(opts))
	//			for _, field := range fields {
	//				val, exists := opts[field.Key]
	//				assert.True(t, exists)
	//				assert.Equal(t, field.val, val)
	//			}
	//		})
	//		t.Run("Replacing fields", func(t *testing.T) {
	//			replacementField := OptionField(fields[0].Key, math.Pi)
	//			optsReplaced := opts.With(replacementField)
	//			assert.Equal(t, len(fields), len(optsReplaced))
	//			val, exists := optsReplaced[replacementField.Key]
	//			assert.True(t, exists)
	//			assert.Equal(t, replacementField.val, val)
	//		})
	//	})
	//	t.Run("Without()", func(t *testing.T) {
	//		assert.Equal(t, 0, len(opts.Without(maps.Keys(opts)...)), "removing all keys should result in a 0-length options set")
	//		optsLessA := opts.Without(keyA)
	//		assert.Equal(t, len(opts)-1, len(optsLessA))
	//		assert.NotContains(t, maps.Keys(optsLessA), keyA, "options without key A should not contain key A")
	//		assert.Contains(t, maps.Keys(optsLessA), keyB, "options without key A should still contain key B")
	//
	//	})
	//})
}
