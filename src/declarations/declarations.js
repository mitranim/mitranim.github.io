/**
 * Globals provided with gulp-wrap.
 */
declare var angular : any;
declare var _ : any;

/**
 * Angular services.
 */
type $Q = {
  all(values: any): Promise;
  reject(value: any): Promise;
  when(value: any): Promise;
}
