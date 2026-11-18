
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

export interface ValidationMethod {
  Method: "DNS-01" | "HTTP-01";
  Confirmed: boolean; // false when inferred (certificates don't record their challenge type)
  Reason: string;
}

export interface SSLDetails {
  HandshakeComplete: boolean;
  DidResume: boolean;
  CipherSuite: number;
  PeerCertificates: CertificateDetails[];
  TLSProtocols: TLSProtocolSupport[];
  HTTPVersions: HTTPVersionSupport[];
  Validation: ValidationMethod | null;
}
