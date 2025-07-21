class ContainerSelector < Formula
  desc "Interactive Docker container selector and command executor"
  homepage "https://github.com/IgorSakharov/container-selector"  # Fixed typo
  url "https://github.com/IgorSakharov/container-selector/archive/refs/tags/v1.0.0.tar.gz"  # Fixed typo
  sha256 "0019dfc4b32d63c1392aa264aed2253c1e0c2fb09216f8e2cc269bbfb8bb49b5"
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
