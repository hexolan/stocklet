import createClient from "openapi-fetch";

import type { paths } from "./schema.d";
import { API_URL } from "../config";

const client = createClient<paths>({ baseUrl: API_URL });

export default client;