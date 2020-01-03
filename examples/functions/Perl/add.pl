sub fx {
  my $ctx = shift;
  my $a = $ctx->req->json->{"a"};
  my $b = $ctx->req->json->{"b"};
  return int($a) + int($b)
}

1;
