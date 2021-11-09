/**
* Copyright 2021 The Nakama Authors
*
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You may obtain a copy of the License at
*
* http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and
* limitations under the License.
*/

using System.Runtime.Serialization;

namespace NakamaSync
{
    /// <summary>
    /// A handshake response from a host to a guest.
    /// </summary>
    internal class HandshakeResponse<T>
    {
        /// <summary>
        /// The store of synced values for a local client to merge in upon
        /// receiving the handshake response.
        /// </summary>
        [DataMember(Name="currentStore"), Preserve]
        public Envelope<T> Store { get; set; }

        /// <summary>
        /// Whether or not the handshake was successful.
        /// </summary>
        [DataMember(Name="success"), Preserve]
        public bool Success { get; set; }

        /// <summary>
        /// Create a handshake response from a host to a guest.
        /// </summary>
        internal HandshakeResponse(Envelope<T> currentState, bool success)
        {
            Store = currentState;
            Success = success;
        }
    }
}
