# UseRepositoryExample

An example of using a generic repository with GORM and OpenTelemetry.



## Structure

```go
type UserExample struct {
    ID    int64  `gorm:"primary_key"`
    Name  string `gorm:"not null"`
    Email string `gorm:"not null"`
}

type UseRepositoryExample interface {
    repository.Repository[UserExample]
    FindEmail(ctx context.Context, email string) (string, error)
}

type useRepositoryExample struct {
    repository.Repository[UserExample]
    db     *gorm.DB
    tracer trace.Tracer
}

func NewUseRepositoryExample(db *gorm.DB) UseRepositoryExample {
    tracer := otel.Tracer("use-repository-example")
    generic := repository.NewRepository[UserExample](db, tracer)
    return &useRepositoryExample{
        Repository: generic,
        db:         db,
        tracer:     tracer,
    }
}

func (u *useRepositoryExample) FindEmail(ctx context.Context, email string) (string, error) {
    ctx, span := u.tracer.Start(ctx, "example.FindEmail")
    defer span.End()

    // TODO: 
    total := "test@test"
    return total, nil
}

```
## Usage

```go
db, _ := gorm.Open(...) //// Configure the database connection

repo := NewUseRepositoryExample(db)

email, err := repo.FindEmail(context.Background(), "user@example.com")
if err != nil {
    log.Fatal(err)
}

fmt.Println("Email:", email)
```