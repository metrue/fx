
struct Input
    a::Number
    b::Number
end

fx = function(input::Input)
    return input.a + input.b
end
