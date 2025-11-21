import type { SSLDetails } from "./certificate";


export interface DomainDetails {
  Domain: string;
  Address: string;
  Resolves: boolean;
}


export interface EndpointDetails {
  Domain: DomainDetails;
  SSL: SSLDetails;
  IsExpired: boolean;
}
