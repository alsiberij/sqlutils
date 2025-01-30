// Package scan has useful functions for handling rows from Query, QueryRow, Exec:
//   - CollectRows - collects slice of items
//   - CollectRowsKV - collects map of items
//   - CollectRow - collects single item
//
// These functions require special row handler row (RowColsCollector, RowColsCollectorKV or RowCollector) that
// accepts row and should return required type. Basic generic implementation are:
//   - Direct scanning: DirectCollectorRowCols, DirectCollectorRowColsKV, DirectCollectorRow
//   - Struct fields by position: StructPosCollectorRowCols, StructPosCollectorRow, StructPosCollectorRowColsKV
//   - Struct fields by tag: StructTagCollectorRowCols, StructTagCollectorRowColsKV
//
// Examples:
//
// DirectCollectorRowCols:
//
//	var ids []int64
//	rows, _ = db.Query("SELECT id FROM users")
//	ids, _ = CollectRows(rows, DirectCollectorRowCols[int64])
//
// DirectCollectorRowColsKV:
//
//	 var users map[int]string
//		rows, _ = db.Query("SELECT id, name FROM users")
//		users, _ = CollectRowsKV(rows, DirectCollectorRowColsKV[int64, string])
//
// DirectCollectorRow:
//
//	 var name string
//		var ok bool
//		row = db.QueryRow("SELECT name FROM users where id = $1")
//		name, ok, _ = CollectRow(row, DirectCollectorRow[string])
//
// StructTagCollectorRowCols:
//
//	 type User struct {
//		    Id    int64  `scan:"id"`
//		    Name  string `scan:"name"`
//		    Email string `scan:"email"`
//	 }
//		var users []User
//		rows, _ = db.Query("SELECT id, name, email FROM users")
//		users, _ = CollectRows(rows, StructTagCollectorRowCols[User])
//
// StructTagCollectorRowColsKV:
//
//	 type User struct {
//		    Name  string `scan:"name"`
//		    Email string `scan:"email"`
//		}
//		var users map[int64]User
//		rows, _ = db.Query("SELECT id, name, email FROM users")
//		users, _ = CollectRowsKV(rows, StructTagCollectorRowColsKV[int64, User])
//
// And so on ...
//
// Package also has Slice for scanning SQL arrays. Example:
//
//	 var arr []int
//		db.QueryRow(`SELECT ARRAY[1, 2, 3]`).Scan(scan.NewSlice(&arr))
//
// Now arr is [1, 2, 3]
package scan
