/*

Packer takes source codes of a function, and pack them into a containerized service, that means there is Dockerfile generated in the output directory

e.g.

	Pack(output, "hello.js") 							# a single file function
	Pack(output, "hello.js", "helper.js") # multiple files function
	Pack(output, "./func/")								# a directory of function
	Pack(output, "hello.js", "./func/")		# a directory and files of function

*/

package packer
