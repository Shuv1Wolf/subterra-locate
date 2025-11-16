import { GEO_HOST } from '../config';

const withOrgId = (url) => {
  const orgId = localStorage.getItem("selectedOrgId");
  if (!orgId) {
    return url;
  }
  const separator = url.includes('?') ? '&' : '?';
  return `${url}${separator}org_id=${orgId}`;
};

const getHeaders = (hasBody = false) => {
  const headers = {};
  if (hasBody) {
    headers["Content-Type"] = "application/json";
  }
  return headers;
};

export const apiClient = {
  get: async (url) => {
    const response = await fetch(withOrgId(url), {
      method: "GET",
      headers: getHeaders(),
    });
    if (!response.ok) {
      const errorData = await response.json().catch(() => ({ message: response.statusText }));
      throw new Error(errorData.message || 'Request failed');
    }
    return response.json();
  },
  post: async (url, data) => {
    const response = await fetch(withOrgId(url), {
      method: "POST",
      headers: getHeaders(true),
      body: JSON.stringify(data),
    });
    if (!response.ok) {
      const errorData = await response.json().catch(() => ({ message: response.statusText }));
      throw new Error(errorData.message || 'Request failed');
    }
    return response.json();
  },
  put: async (url, data) => {
    const response = await fetch(withOrgId(url), {
      method: "PUT",
      headers: getHeaders(true),
      body: JSON.stringify(data),
    });
    if (!response.ok) {
      const errorData = await response.json().catch(() => ({ message: response.statusText }));
      throw new Error(errorData.message || 'Request failed');
    }
    return response.json();
  },
  delete: async (url) => {
    const response = await fetch(withOrgId(url), {
      method: "DELETE",
      headers: getHeaders(),
    });
    if (!response.ok) {
      const errorData = await response.json().catch(() => ({ message: response.statusText }));
      throw new Error(errorData.message || 'Request failed');
    }
    const text = await response.text();
    return text ? JSON.parse(text) : Promise.resolve();
  },
};

export const getDeviceHistory = async (mapId, deviceId, from, to, take, skip) => {
  const orgId = localStorage.getItem("selectedOrgId");
  let url = `${GEO_HOST}/api/v1/geo/history?map_id=${mapId}&from=${from}&to=${to}&take=${take}&total=true&skip=${skip}&entity_id=${deviceId}`;
  if (orgId) {
    url += `&org_id=${orgId}`;
  }
  
  // Using fetch directly to ensure org_id is correctly handled without side effects from apiClient wrapper
  const response = await fetch(url, {
    method: "GET",
    headers: getHeaders(),
  });

  if (!response.ok) {
    const errorData = await response.json().catch(() => ({ message: response.statusText }));
    throw new Error(errorData.message || 'Request failed');
  }
  return response.json();
};

export const getMaps = async (orgId, skip = 0, take = 10) => {
  return apiClient.get(
    `${GEO_HOST}/api/v1/geo/map?org_id=${orgId}&total=true&skip=${skip}&take=${take}`
  );
};

export const getMap = async (mapId) => {
  return apiClient.get(`${GEO_HOST}/api/v1/geo/map/${mapId}`);
};

export const deleteMap = async (mapId) => {
  return apiClient.delete(`${GEO_HOST}/api/v1/geo/map/${mapId}`);
};

export const createMap = async (mapData) => {
  return apiClient.post(`${GEO_HOST}/api/v1/geo/map`, mapData);
};

export const updateMap = async (mapData) => {
  return apiClient.put(`${GEO_HOST}/api/v1/geo/map`, mapData);
};

export const uploadMap = async (file, mapId) => {
  const formData = new FormData();
  formData.append('file', file);

  const url = withOrgId(`${GEO_HOST}/api/v1/geo/map/upload?id=${mapId}`);

  const response = await fetch(url, {
    method: 'POST',
    body: formData,
  });

  if (!response.ok) {
    const errorData = await response.json().catch(() => ({ message: response.statusText }));
    throw new Error(errorData.message || 'Request failed');
  }

  return response.json();
};
