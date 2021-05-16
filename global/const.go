package global

const (
	dburl       = "mongodb+srv://Hitesh1103:mzlRpnSLJmtFHCss@practicecluster.7ie7c.mongodb.net/myFirstDatabase?retryWrites=true&w=majority"
	dbname      = "blog-application"
	performance = 100
)

var (
	jwtSecret = []byte("blogsecret")
)

const (
	EmailRegex = "^[a-zA-Z0-9-.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$"
)
