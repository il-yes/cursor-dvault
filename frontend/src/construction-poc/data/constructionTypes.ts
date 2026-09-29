/**
 * View-model types for the Construction Functional Data Boundary.
 *
 * These are the contracts the presentation layer already depends on. They are
 * intentionally expressed in UI terms, not as a second domain model, and are
 * shared by every ConstructionDataProvider implementation.
 */

export interface CreateProjectParams {
  code: string;
  name: string;
  type: string;
  location: string;
}

export interface ProcurementData {
  id: string;
  requirementId: string;
  projectId: string;
  projectCode: string;
  projectName: string;
  code: string;
  materialId: string;
  materialName: string;
  specification: string;
  quantity: number;
  unit: string;
  targetPhase: string;
  status: string;
  requiredDate: string;
  priority: string;
  siteId: string;
  siteName: string;
  invitedSuppliersCount: number;
  offersReceivedCount: number;
  offerId: string;
  offerReference: string;
  supplierId: string;
  supplierName: string;
  totalPrice: string;
  unitPrice: string;
  promisedDeliveryDate: string;
  offerStatus: string;
  isVerifiedSupplier: boolean;
}

export interface LogisticsOverviewData {
  delivery: {
    id: string;
    reference: string;
    status: string;
    plannedDeliveryDate: string;
    eta: string;
    actualDeliveryDate?: string;
    quantity: number;
    unit: string;
    deliveryNotes?: string;
  };
  project: {
    id: string;
    code: string;
    name: string;
    type: string;
  };
  requirement: {
    id: string;
    code: string;
    targetPhase: string;
  };
  material: {
    id: string;
    name: string;
    specification: string;
  };
  supplier: {
    id: string;
    name: string;
  };
  offer: {
    id: string;
    reference: string;
    totalPrice: string;
    unitPrice: string;
    promisedDeliveryDate: string;
  };
  site: {
    id: string;
    name: string;
  };
  transport: {
    id: string;
    reference: string;
    vehicle: string;
    driver: string;
    status: string;
    origin: string;
    destination: string;
    route: string;
    plannedDeparture: string;
    actualDeparture: string;
    plannedArrival: string;
    eta: string;
    actualArrival?: string;
    constraints?: string[];
    delayReason?: string;
  };
  issue?: {
    id: string;
    reference: string;
    title: string;
    description: string;
    severity: string;
    status: string;
    reportedBy: string;
    reportedAt: string;
    impact: string;
    evidenceReferences?: string[];
  };
  decision?: {
    id: string;
    reference: string;
    subject: string;
    context: string;
    technicalAssessment?: string;
    risksIdentified?: string[];
    participantsConsulted?: string[];
    optionsConsidered?: string[];
    decision: string;
    decidedBy: string;
    decisionDate: string;
    consequence: string;
    status: string;
  };
  inspection?: {
    id: string;
    reference: string;
    type: string;
    inspector: string;
    inspectionDate: string;
    status: string;
    result: string;
    criteria?: string[];
    findings?: string[];
    notes?: string;
    evidenceReferences?: string[];
  };
}
