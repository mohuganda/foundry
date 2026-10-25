class Foundry < Formula
  desc "Create and extend projects from reusable foundations"
  homepage "https://github.com/mohuganda/foundry"
  version "1.4.0"

  on_macos do
    on_arm do
      url "https://github.com/mohuganda/foundry/releases/download/v1.4.0/foundry_1.4.0_darwin_arm64.tar.gz"
      sha256 "393b9ebf1b89399afd22a4593944837cf08c24402ba761063c341677cb9fe1c8"
    end
    on_intel do
      url "https://github.com/mohuganda/foundry/releases/download/v1.4.0/foundry_1.4.0_darwin_amd64.tar.gz"
      sha256 "667164abb5a577ec1de8d687807b975c4612a0b9ee5478772e1cfa491de0a151"
    end
  end

  on_linux do
    on_arm do
      url "https://github.com/mohuganda/foundry/releases/download/v1.4.0/foundry_1.4.0_linux_arm64.tar.gz"
      sha256 "51e896f73179be557ef860fa7ae6ec7fbcbc43f9de198929fc75095e8e8362b1"
    end
    on_intel do
      url "https://github.com/mohuganda/foundry/releases/download/v1.4.0/foundry_1.4.0_linux_amd64.tar.gz"
      sha256 "6b6b8cdd540515917cb070e73bdc4e1e1cbbf78186c167700a721dc27092ec6c"
    end
  end

  def install
    bin.install "foundry"
    doc.install "README.md", "CHANGELOG.md", "docs/commands.md", "docs/workspaces.md"
  end

  test do
    assert_match "foundry #{version}", shell_output("#{bin}/foundry version")
    system bin/"foundry", "new", "brew-smoke", "--template", "rust", "--no-hooks"
    assert_path_exists testpath/"brew-smoke/Cargo.toml"
  end
end
