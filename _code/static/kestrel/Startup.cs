using System.Threading.Tasks;
using Microsoft.AspNetCore.Builder;
using Microsoft.AspNetCore.Http;
using System.Text;

namespace netcore
{
    public class Startup
    {
        public void Configure(IApplicationBuilder app)
        {
            app.UseMiddleware<PlaintextMiddleware>();
        }
    }

    public class PlaintextMiddleware
    {
        private static Encoding Encoding = Encoding.UTF8;

        public PlaintextMiddleware(RequestDelegate next)
        {
        }

        public Task Invoke(HttpContext httpContext)
        {
            var bytes = Encoding.GetBytes("Index");
            return httpContext.Response.Body.WriteAsync(bytes, 0, bytes.Length);
        }
    }
}
