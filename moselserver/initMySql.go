/*
 * Copyright 2017 Robin Engel
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package moselserver

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

func (server *MoselServer) initMySql(driverName string, dataSourceName string) (SqlDataSource, error) {
	var db *sql.DB
	var err error

	if db, err = sql.Open(driverName, dataSourceName); err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return NewSqlDataSource(driverName, db), nil
}
