# Usecase Layer

Application business workflows live here.

Usecases may depend on domain types and repository interfaces. They must not
know about Gin, HTTP response helpers, GORM models, or Redis clients directly.
