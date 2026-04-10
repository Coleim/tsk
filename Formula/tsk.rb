class Tsk < Formula
  desc "Terminal task manager with TUI"
  homepage "https://github.com/Coleim/tsk"
  version "0.3.0"
  license "MIT"

  on_macos do
    on_intel do
      url "https://github.com/Coleim/tsk/releases/download/v0.3.0/tsk_darwin_amd64.tar.gz"
      sha256 "731b4502f80aedb10caa2ec5f56e46515766c51032552970c45b22bfbfc37116"
    end
    on_arm do
      url "https://github.com/Coleim/tsk/releases/download/v0.3.0/tsk_darwin_arm64.tar.gz"
      sha256 "edfbabaa035317e48b01c96ccd65e4339453631fdc4700225600038dc6f1987d"
    end
  end

  on_linux do
    on_intel do
      url "https://github.com/Coleim/tsk/releases/download/v0.3.0/tsk_linux_amd64.tar.gz"
      sha256 "ff185b8ed5cb14e621867e0a02713fb3280c4e64da51347b1dda42f4a7959c1c"
    end
    on_arm do
      url "https://github.com/Coleim/tsk/releases/download/v0.3.0/tsk_linux_arm64.tar.gz"
      sha256 "fbc4b60575b2b09783bfde5dadd4fc83c9c8616e14e2f82f8885cad2ab4e1aa6"
    end
  end

  def install
    bin.install "tsk"
  end

  test do
    system "\#{bin}/tsk", "--version"
  end
end
