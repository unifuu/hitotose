<?php

namespace App\Http\Controllers;

use Illuminate\Http\Request;

class ApiController extends Controller
{
    public function getCsrf(Request $request)
    {
        // Generate a new CSRF token
        return response()->json(['csrf_token' => csrf_token()]);
    }
}