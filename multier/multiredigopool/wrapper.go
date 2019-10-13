package multiredigopool

import redigo "github.com/gomodule/redigo/redis"

// Get ...
func Get() redigo.Conn {
	s, _ := GetN(d)
	return s
}

// GetN ...
func GetN(name string) (redigo.Conn, error) {
	pool, err := RetrieveN(name)
	if err != nil {
		return nil, err
	}
	return pool.Get(), nil
}

// MustGetN ...
func MustGetN(name string) redigo.Conn {
	s, err := GetN(d)
	if err != nil {
		panic(err)
	}
	return s
}
