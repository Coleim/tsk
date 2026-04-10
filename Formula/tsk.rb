class Tsk < Formula
  desc "Terminal task manager with TUI"
  homepage "https://github.com/Coleim/tsk"
  version "0.3.0"
  license "MIT"

  on_macos do
    on_intel do
      url "https://github.com/Coleim/tsk/releases/download/v0.3.0/tsk_darwin_amd64.tar.gz"
      sha256 "b2590687d4de4c837ace3b0c5d2beb4a387092902ee39fcd1c8bb9561e2e11f4"
    end
    on_arm do
      url "https://github.com/Coleim/tsk/releases/download/v0.3.0/tsk_darwin_arm64.tar.gz"
      sha256 "382d2147b01a90bdb53f0c4a6d01184d0dffc7adf22953d527da3171dda2e0d4"
    end
  end

  on_linux do
    on_intel do
      url "https://github.com/Coleim/tsk/releases/download/v0.3.0/tsk_linux_amd64.tar.gz"
      sha256 "f2b9a79b71057aa5a8bd448b28746206db8d31597b244570a9a76ed3cd179621"
    end
    on_arm do
      url "https://github.com/Coleim/tsk/releases/download/v0.3.0/tsk_linux_arm64.tar.gz"
      sha256 "2cb80c5685b125136ec98020bcad1de07aeec591a2ec1d4942845ba4d632ba24"
    end
  end

  def install
    bin.install "tsk"
  end

  test do
    system "\#{bin}/tsk", "--version"
  end
end
