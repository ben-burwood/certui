
export interface CertificateDetails {
  Subject: string;
  Issuer: string;
  NotBefore: string; // ISO date string
  NotAfter: string;  // ISO date string
  SignatureAlgorithm: string;
  PublicKeyAlgorithm: string;
}

export interface TLSProtocolSupport {
  Protocol: string;
  Supported: boolean;
}

export interface HTTPVersionSupport {
  Version: string;
  Supported: boolean;
}

export interface SSLDetails {
  HandshakeComplete: boolean;
  DidResume: boolean;
  CipherSuite: number;
  PeerCertificates: CertificateDetails[];
  TLSProtocols: TLSProtocolSupport[];
  HTTPVersions: HTTPVersionSupport[];
}
