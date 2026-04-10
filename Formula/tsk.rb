class Tsk < Formula
  desc "Terminal task manager with TUI"
  homepage "https://github.com/Coleim/tsk"
  version "0.3.0"
  license "MIT"

  on_macos do
    on_intel do
      url "https://github.com/Coleim/tsk/releases/download/v0.3.0/tsk_darwin_amd64.tar.gz"
      sha256 "98277149d7a3ab2579f102824744db31d48e1d9bb04c4db2df2eb0893567a388"
    end
    on_arm do
      url "https://github.com/Coleim/tsk/releases/download/v0.3.0/tsk_darwin_arm64.tar.gz"
      sha256 "4c0a8a756595673e6b9deaf07a1bd9a9cca5e981e1ac7a9f02e0a1627be7072a"
    end
  end

  on_linux do
    on_intel do
      url "https://github.com/Coleim/tsk/releases/download/v0.3.0/tsk_linux_amd64.tar.gz"
      sha256 "74523e3b15bb059540d21277367b949d7dcb7c4de5fad647db260f0cd30466c2"
    end
    on_arm do
      url "https://github.com/Coleim/tsk/releases/download/v0.3.0/tsk_linux_arm64.tar.gz"
      sha256 "2ab4b0b3a08570a0cc4601532a8d0698fc4d08e696f9e5e59e36188762b4b320"
    end
  end

  def install
    bin.install "tsk"
  end

  test do
    system "\#{bin}/tsk", "--version"
  end
end
