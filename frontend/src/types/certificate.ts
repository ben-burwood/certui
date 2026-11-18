
export interface CertificateDetails {
  Subject: string;
  Issuer: string;
  NotBefore: string; // ISO date string
  NotAfter: string;  // ISO date string
  SignatureAlgorithm: string;
  PublicKeyAlgorithm: string;
  DNSNames: string[] | null;
  Fingerprint: string; // SHA-256, hex-encoded
  IsWildcard: boolean;
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

export interface CAAIssuer {
  CA: string; // "" means no CA may issue
  ValidationMethods: string[] | null; // null means unrestricted
  AccountURI: string;
}

export interface CAAPolicy {
  Domain: string; // the host itself or the parent the records were inherited from
  Issuers: CAAIssuer[] | null;
  WildcardIssuers: CAAIssuer[] | null;
}

export interface SSLDetails {
  HandshakeComplete: boolean;
  DidResume: boolean;
  CipherSuite: number;
  PeerCertificates: CertificateDetails[];
  TLSProtocols: TLSProtocolSupport[];
  HTTPVersions: HTTPVersionSupport[];
  Validation: ValidationMethod | null;
  CAA: CAAPolicy | null; // null when no CAA records exist (any CA may issue) or the lookup failed
  CAAError: string; // set when the CAA lookup failed
}
