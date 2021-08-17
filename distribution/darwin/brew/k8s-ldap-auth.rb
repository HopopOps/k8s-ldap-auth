class K8sLdapAuth < Formula
  desc "Kubernetes webhook token authentication plugin implementation using ldap"
  homepage "https://github.com/vbouchaud/k8s-ldap-auth/"
  url "https://github.com/vbouchaud/k8s-ldap-auth/archive/refs/tags/v3.0.0.tar.gz"
  sha256 "bd46ebae08fc850065db7ab2ed38c47dbf85156585c31323d3157101f42e35da"
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
