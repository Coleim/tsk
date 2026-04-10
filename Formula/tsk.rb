class Tsk < Formula
  desc "Terminal task manager with TUI"
  homepage "https://github.com/Coleim/tsk"
  version "0.1.0"
  license "MIT"

  on_macos do
    on_intel do
      url "https://github.com/Coleim/tsk/releases/download/v0.1.0/tsk_darwin_amd64.tar.gz"
      sha256 ""
    end
    on_arm do
      url "https://github.com/Coleim/tsk/releases/download/v0.1.0/tsk_darwin_arm64.tar.gz"
      sha256 ""
    end
  end

  on_linux do
    on_intel do
      url "https://github.com/Coleim/tsk/releases/download/v0.1.0/tsk_linux_amd64.tar.gz"
      sha256 ""
    end
    on_arm do
      url "https://github.com/Coleim/tsk/releases/download/v0.1.0/tsk_linux_arm64.tar.gz"
      sha256 ""
    end
  end

  def install
    bin.install "tsk"
  end

  test do
    system "\#{bin}/tsk", "--version"
  end
end
