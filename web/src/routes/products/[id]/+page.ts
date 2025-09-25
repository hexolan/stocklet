// import { error } from '@sveltejs/kit';
import type { PageLoad } from './$types';

export const load: PageLoad = ({ params }) => {
  return {
    id: params.id
  };

  // error(404, 'Product not found');
};