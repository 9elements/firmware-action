package namespace.exceptions

import data.namespaces

exception[ns] {
	# Ignore: Image user should not be 'root'
	ns := data.namespaces[_]
	startswith(ns, "builtin.dockerfile.DS002")
}

exception[ns] {
	# Ignore: No HEALTHCHECK defined
	ns := data.namespaces[_]
	startswith(ns, "builtin.dockerfile.DS026")
}

exception[ns] {
	# Ignore: 'RUN cd ...' to change directory
	ns := data.namespaces[_]
	startswith(ns, "builtin.dockerfile.DS013")
}
