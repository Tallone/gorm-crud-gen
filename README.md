# GORM CRUD Generator

This CLI tool generates CRUD services and Gin handlers with Swaggo comments from GORM structs.

## Installation

```bash
go get github.com/Tallone/gorm-crud-gen
```

## Usage

```bash
gorm-crud-gen generate [input_file] -o [output_directory] -p [package_name]
```

- `input_file`: Path to the file containing the GORM struct
- `-o, --output`: Output directory for generated files (default: current directory)
- `-p, --package`: Package name for generated files (default: "main")

## Example

Given a file `user.go` with the following content:

```go
package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name  string `gorm:"index"`
	Email string `gorm:"uniqueIndex"`
	Age   int
}
```

Run the following command:

```bash
gorm-crud-gen generate user.go -o generated -p myapp
```

This will generate CRUD services and handlers for the `User` struct in the `generated` directory.

## Generated Files

- `services/user_service.go`: Contains CRUD operations for the User model
- `handlers/user_handler.go`: Contains Gin handlers with Swaggo comments for the User API

## Dependencies


## License

This project is licensed under the MIT License.