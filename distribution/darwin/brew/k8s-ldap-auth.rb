class K8sLdapAuth < Formula
  desc "Kubernetes webhook token authentication plugin implementation using ldap"
  homepage "https://github.com/vbouchaud/k8s-ldap-auth/"
  url "https://github.com/vbouchaud/k8s-ldap-auth/archive/refs/tags/v2.0.1.tar.gz"
  sha256 "2b1dc9fe80ffa06593d981722e644c7a59805138349e08f824775200096bb58a"
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
