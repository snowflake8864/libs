package public

import ()

type Void interface{}
type Dtime uint64

var Error = []string{
	"ERR_NOERR",
	"ERR_EPERM",
	"ERR_ENOENT",
	"ERR_EINTR",
	"ERR_EIO",
	"ERR_ENXIO",
	"ERR_EAGAIN",
	"ERR_ENOMEM",
	"ERR_EACCESS",
	"ERR_EFAULT",
	"ERR_EBUSY",
	"ERR_EEXIST",
	"ERR_EINVAL",
	"ERR_ENONET",
	"ERR_EPROTO",
	"ERR_ENOPROTOOPT",
	"ERR_EPROTONOSUPPORT",
	"ERR_EOPNOTSUPP",
	"ERR_EADDRINUSE",
	"ERR_EADDRNOTAVAIL",
	"ERR_ENETDOWN",
	"ERR_ENETUNREACH",
	"ERR_ECONNRESET",
	"ERR_EISCONN",
	"ERR_ENOTCONN",
	"ERR_ESHUTDOWN",
	"ERR_ETIMEDOUT",
	"ERR_ECONNREFUSED",
	"ERR_EHOSTDOWN",
	"ERR_EHOSTUNREACH",
}

func Perror(t int) string {
	return Error[t]
}

type Derror interface {
	//Perror()
}

const (
	ERR_NOERR = iota
	ERR_EPERM
	ERR_ENOENT
	ERR_EINTR
	ERR_EIO
	ERR_ENXIO
	ERR_EAGAIN
	ERR_ENOMEM
	ERR_EACCESS
	ERR_EFAULT
	ERR_EBUSY
	ERR_EEXIST
	ERR_EINVAL
	ERR_ENONET
	ERR_EPROTO
	ERR_ENOPROTOOPT
	ERR_EPROTONOSUPPORT
	ERR_EOPNOTSUPP
	ERR_EADDRINUSE
	ERR_EADDRNOTAVAIL
	ERR_ENETDOWN
	ERR_ENETUNREACH
	ERR_ECONNRESET
	ERR_EISCONN
	ERR_ENOTCONN
	ERR_ESHUTDOWN
	ERR_ETIMEDOUT
	ERR_ECONNREFUSED
	ERR_EHOSTDOWN
	ERR_EHOSTUNREACH
)
