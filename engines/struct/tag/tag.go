package tag

type Tag[F any] struct {
	Name   string
	Format *F
}
