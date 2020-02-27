package erro

type Err interface {
	Erro() (error, string)
}

type WrapError struct {
	Action string
	Err    error
}

type ConfigError struct {
	Action string
	Err    error
}

type DBConnError struct {
	Action string
	Err    error
}

type AtoiError struct {
	Action string
	Err    error
}

type DateError struct {
	Action string
	Err    error
}

func (e *WrapError) Erro() (error, string) {
	return e.Err, e.Action
}

func (e *ConfigError) Erro() (error, string) {
	return e.Err, e.Action
}

func (e *DBConnError) Erro() (error, string) {
	return e.Err, e.Action
}

func (e *AtoiError) Erro() (error, string) {
	return e.Err, e.Action
}

func NewWrapError(action string, err error) *WrapError {
	return &WrapError{action, err}
}

func NewConfigError(action string, err error) *ConfigError {
	return &ConfigError{action, err}
}

func NewDBConnError(action string, err error) *DBConnError {
	return &DBConnError{action, err}
}

func NewAtoiError(action string, err error) *AtoiError {
	return &AtoiError{action, err}
}
