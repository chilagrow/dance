// Copyright 2021 FerretDB Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package diff

import (
	"math"
	"testing"

	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func TestFloatValues(t *testing.T) {
	t.Parallel()

	ctx, db := setup(t)

	t.Run("Insert", func(t *testing.T) {
		t.Parallel()

		for name, tc := range map[string]struct {
			doc      bson.D
			expected mongo.CommandError
		}{
			"NaN": {
				doc: bson.D{{"_id", "1"}, {"foo", math.NaN()}},
				expected: mongo.CommandError{
					Code: 2,
					Name: "BadValue",
					Message: `wire.OpMsg.Document: validation failed for { insert: "insert-NaN", ordered: true, ` +
						`$db: "testfloatvalues", documents: [ { _id: "1", foo: nan.0 } ] } with: NaN is not supported`,
				},
			},
			"NegativeZero": {
				doc: bson.D{{"_id", "1"}, {"foo", math.Copysign(0.0, -1)}},
				expected: mongo.CommandError{
					Code: 2,
					Name: "BadValue",
					Message: `wire.OpMsg.Document: validation failed for { insert: "insert-NegativeZero", ordered: true, ` +
						`$db: "testfloatvalues", documents: [ { _id: "1", foo: -0.0 } ] } with: -0 is not supported`,
				},
			},
		} {
			name, tc := name, tc

			t.Run(name, func(t *testing.T) {
				t.Parallel()

				_, err := db.Collection("insert-"+name).InsertOne(ctx, tc.doc)

				t.Run("FerretDB", func(t *testing.T) {
					t.Parallel()

					assertEqualError(t, tc.expected, err)
				})

				t.Run("MongoDB", func(t *testing.T) {
					t.Parallel()

					require.NoError(t, err)
				})
			})
		}
	})

	// TODO https://github.com/FerretDB/dance/issues/266
	/*t.Run("Update", func(t *testing.T) {

	})*/

	// TODO https://github.com/FerretDB/dance/issues/266
	/*t.Run("FindAndModify", func(t *testing.T) {

	})*/

	t.Run("UpdateOne", func(t *testing.T) {
		t.Parallel()

		for name, tc := range map[string]struct {
			filter      bson.D
			insert      bson.D
			update      bson.D
			expectedErr mongo.CommandError
		}{
			"MulMaxFloat64": {
				filter: bson.D{{"_id", "number"}},
				insert: bson.D{{"_id", "number"}, {"v", int32(42)}},
				update: bson.D{{"$mul", bson.D{{"v", math.MaxFloat64}}}},
				expectedErr: mongo.CommandError{
					Code:    2,
					Name:    "BadValue",
					Message: `invalid value: { "v": +Inf } (infinity values are not allowed)`,
				},
			},
		} {
			name, tc := name, tc

			t.Run(name, func(t *testing.T) {
				t.Parallel()

				collection := db.Collection("update-one-" + name)
				_, err := collection.InsertOne(ctx, tc.insert)
				require.NoError(t, err)

				_, err = collection.UpdateOne(ctx, tc.filter, tc.update)

				t.Run("FerretDB", func(t *testing.T) {
					t.Parallel()

					assertEqualError(t, tc.expected, err)
				})

				t.Run("MongoDB", func(t *testing.T) {
					t.Parallel()

					require.NoError(t, err)
				})
			})
		}
	})

	t.Run("UpdateOneMul", func(t *testing.T) {
		t.Parallel()

		filter := bson.D{{"_id", "number1"}}

		collection := db.Collection("update-one-mul")
		_, err := collection.InsertOne(ctx, bson.D{{"_id", "number1"}, {"v", float64(-123456789)}})
		require.NoError(t, err)

		_, err = collection.UpdateOne(ctx, filter, bson.D{{"$mul", bson.D{{"v", float64(0)}}}})
		require.NoError(t, err)

		var updated bson.D
		err = collection.FindOne(ctx, filter).Decode(&updated)
		require.NoError(t, err)

		expected := bson.D{{"_id", "number1"}, {"v", math.Copysign(0.0, -1)}}

		t.Run("FerretDB", func(t *testing.T) {
			t.Parallel()

			require.Equal(t, expected, updated)
		})

		t.Run("MongoDB", func(t *testing.T) {
			t.Parallel()

			require.Equal(t, expected, updated)
		})
	})
}
