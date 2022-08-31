export const headers = {
  'Accept': 'application/json',
  'Content-Type': 'application/json',
};

export const options = {
  headers,
  credentials: 'include',
};


export const PROTOCOL = 'https://';
export const DOMAIN = process.env.NODE_ENV === 'development' ? 'localhost' : location.hostname;
export const FRONT_DOMAIN = process.env.NODE_ENV === 'development' ? 'localhost:4321' : location.hostname;
export const BASE_URL = PROTOCOL + DOMAIN;
