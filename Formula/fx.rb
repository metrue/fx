class Fx < Formula
  desc "fx, Poor man's Function as a Service framework"
  homepage "https://github.com/metrue/fx"
  url "https://github.com/metrue/fx/releases/download/v0.0.0/fx_0.0.0_macOS_64-bit.tar.gz"
  version "0.0.0"
  sha256 "5cd775ba46a2a8f2220eaddd78f411a85ad5fd4c72c171d44cb65846823ce74f"

  depends_on "git"depends_on "zsh"

  def install
    bin.install "fx"
  end

  def caveats
    "fx help to get know more"
  end

  test do
    
  end
end
