open("/app/REQUIRE") do f
	deps = readlines(f)
	for d in deps
		Pkg.add(d)
	end
end
Pkg.build("HttpParser")
