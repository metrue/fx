using HttpServer
import JSON
import Unmarshal

include("fx.jl")

srv = Server() do req::Request, res::Response
    if length(req.data) == 0
        res.status = 400
    else
        parsed_data = JSON.Parser.parse(join([Char(v) for v in req.data]))
        data = Unmarshal.unmarshal(Input, parsed_data)
        string(fx(data))
    end
end

run(srv, 3000)
