import std.json;

long executeFx(JSONValue input)
{
    return input["a"].integer + input["b"].integer;
}