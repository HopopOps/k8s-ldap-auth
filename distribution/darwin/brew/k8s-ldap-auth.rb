class K8sLdapAuth < Formula
  desc "Kubernetes webhook token authentication plugin implementation using ldap"
  homepage "https://github.com/vbouchaud/k8s-ldap-auth/"
  url "https://github.com/vbouchaud/k8s-ldap-auth/archive/refs/tags/v3.2.1.tar.gz"
  sha256 "2517c85e6c6e0aebd0062e3e5391511f7b50d6d4237796426e2f1ffb17a7b94d"
  license "MPL-2.0"

  depends_on "go" => :build
  depends_on "gnu-sed" => :build

  def install
    ENV["VERSION"] = "#{version}"
    ENV["SED"] = "gsed"

    system "make", "k8s-ldap-auth"
    system "mkdir", "-p", "#{prefix}/bin/"
    system "cp", "k8s-ldap-auth", "#{prefix}/bin/"
  end

  test do
    system "false"
  end
end
