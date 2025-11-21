
export interface CertificateDetails {
  Subject: string;
  Issuer: string;
  NotBefore: string; // ISO date string
  NotAfter: string;  // ISO date string
  SignatureAlgorithm: string;
  PublicKeyAlgorithm: string;
}

export interface SSLDetails {
  Version: number;
  HandshakeComplete: boolean;
  DidResume: boolean;
  CipherSuite: number;
  PeerCertificates: CertificateDetails[];
}
