class Foundry < Formula
  desc "Create and extend projects from reusable foundations"
  homepage "https://github.com/mohuganda/foundry"
  version "1.3.0"

  on_macos do
    on_arm do
      url "https://github.com/mohuganda/foundry/releases/download/v1.3.0/foundry_1.3.0_darwin_arm64.tar.gz"
      sha256 "66b4ee1bda6fbc1f4aa78529f2ab0d1d24ce49c3748180f753d3cd2b25ab048b"
    end
    on_intel do
      url "https://github.com/mohuganda/foundry/releases/download/v1.3.0/foundry_1.3.0_darwin_amd64.tar.gz"
      sha256 "7084a8760ab09a7c0d56afbd30862e09a4deba15e706d5a21adde277192d8886"
    end
  end

  on_linux do
    on_arm do
      url "https://github.com/mohuganda/foundry/releases/download/v1.3.0/foundry_1.3.0_linux_arm64.tar.gz"
      sha256 "1a7ea9791e800a9a75334cad036a8cbed3e25fd56748ce34d38f10f4a95e113e"
    end
    on_intel do
      url "https://github.com/mohuganda/foundry/releases/download/v1.3.0/foundry_1.3.0_linux_amd64.tar.gz"
      sha256 "396b9211215abc1467daeba4ebd20159009845de72b7b61fa68c29c1d6d6cf59"
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
