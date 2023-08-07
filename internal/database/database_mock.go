/*
 * Copyright 2020-2023 Luke Whritenour

 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at

 *     http://www.apache.org/licenses/LICENSE-2.0

 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package database

import (
	"context"
	"database/sql"

	_ "github.com/lib/pq"
)

type Mock struct {
	*sql.DB
}

func NewMock() (Database, error) {
	return &Mock{&sql.DB{}}, nil
}

func (m *Mock) Migrate(ctx context.Context) error {
	return nil
}

func (m *Mock) GetDocument(ctx context.Context, id string) (Document, error) {
	return Document{}, nil
}

func (m *Mock) CreateDocument(ctx context.Context, id, content string) error {
	return nil
}
