class K8sLdapAuth < Formula
  desc "Kubernetes webhook token authentication plugin implementation using ldap"
  homepage "https://github.com/vbouchaud/k8s-ldap-auth/"
  url "https://github.com/vbouchaud/k8s-ldap-auth/archive/refs/tags/v3.1.0.tar.gz"
  sha256 "3ad203d70ac8ed1be0b5806a527d708da084807b5fc8d496c07d7bb9c99c0298"
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
