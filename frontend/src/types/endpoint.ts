import type { SSLDetails } from "./certificate";
import type { WhoisDetails } from "./whois";


export interface DomainDetails {
  Domain: string;
  Address: string;
  Resolves: boolean;
}


export interface EndpointDetails {
  Domain: DomainDetails;
  Whois: WhoisDetails;
  SSL: SSLDetails | null;
}
