class K8sLdapAuth < Formula
  desc "Kubernetes webhook token authentication plugin implementation using ldap"
  homepage "https://github.com/vbouchaud/k8s-ldap-auth/"
  url "https://github.com/vbouchaud/k8s-ldap-auth/archive/refs/tags/v2.0.0.tar.gz"
  sha256 "9397ad92d6910b922cb501ef02b52a2b20f6a4f1476f62500fbdd29dae2031c6"
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
