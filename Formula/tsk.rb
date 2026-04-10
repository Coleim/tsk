class Tsk < Formula
  desc "Terminal task manager with TUI"
  homepage "https://github.com/Coleim/tsk"
  version "0.3.1"
  license "MIT"

  on_macos do
    on_intel do
      url "https://github.com/Coleim/tsk/releases/download/v0.3.1/tsk_darwin_amd64.tar.gz"
      sha256 "6ad297fa8f97b2db94334aef938770a465b3466291a4b5a5f83716c026b34af9"
    end
    on_arm do
      url "https://github.com/Coleim/tsk/releases/download/v0.3.1/tsk_darwin_arm64.tar.gz"
      sha256 "9bce25e847aa380f3d423c18f36d70c82a902190c114826f99fa2765b7b93834"
    end
  end

  on_linux do
    on_intel do
      url "https://github.com/Coleim/tsk/releases/download/v0.3.1/tsk_linux_amd64.tar.gz"
      sha256 "8c481288c5ad0a22a3642eb85223ad0b9a4af662a26eea2eb3533a9deba14a63"
    end
    on_arm do
      url "https://github.com/Coleim/tsk/releases/download/v0.3.1/tsk_linux_arm64.tar.gz"
      sha256 "9525d4c0593b4b5ec3a820850347a1a0b391806be26421904068e617aa7a8dc8"
    end
  end

  def install
    bin.install "tsk"
  end

  test do
    system "\#{bin}/tsk", "--version"
  end
end
