class Tsk < Formula
  desc "Terminal task manager with TUI"
  homepage "https://github.com/Coleim/tsk"
  version "0.2.0"
  license "MIT"

  on_macos do
    on_intel do
      url "https://github.com/Coleim/tsk/releases/download/v0.2.0/tsk_darwin_amd64.tar.gz"
      sha256 "f680faddde8981b48dc37c354db0b97ec8a38a3df2c23b8bc6668999f2333341"
    end
    on_arm do
      url "https://github.com/Coleim/tsk/releases/download/v0.2.0/tsk_darwin_arm64.tar.gz"
      sha256 "34bec1f1395c3acd12235d8841fc7b798796b900b27980aa650c7e9fb15ee426"
    end
  end

  on_linux do
    on_intel do
      url "https://github.com/Coleim/tsk/releases/download/v0.2.0/tsk_linux_amd64.tar.gz"
      sha256 "211c2c3a4ef75961356a77123922db6049802b925fd66b2e35036000e7f88c4d"
    end
    on_arm do
      url "https://github.com/Coleim/tsk/releases/download/v0.2.0/tsk_linux_arm64.tar.gz"
      sha256 "909e90af2b45bf6ad1d02098c1868c2b077c9a02d54c06fb31a6599fd5c7644f"
    end
  end

  def install
    bin.install "tsk"
  end

  test do
    system "\#{bin}/tsk", "--version"
  end
end
