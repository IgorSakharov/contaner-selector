class ContainerSelector < Formula
  desc "Interactive Docker container selector and command executor"
  homepage "https://github.com/IgorSakharov/container-selector"  # Fixed typo
  url "https://github.com/IgorSakharov/container-selector/archive/refs/tags/v1.0.0.tar.gz"  # Fixed typo
  sha256 "a396b37095023498a4e4ff36006dd5d0c6ce72af947a228dc90e727fbbc384df"
  license "MIT"

  depends_on "go" => :build

  def install
    system "go", "build", *std_go_args(output: bin/"container-selector")
  end

  test do
    # Test that the binary exists and shows help (since --version might not be implemented)
    system "#{bin}/container-selector", "--help"
  end
end
