package error

type CachesError string

func (e CachesError) Error() string {
	return string(e)
}

const Nil = CachesError("caches: nil")
